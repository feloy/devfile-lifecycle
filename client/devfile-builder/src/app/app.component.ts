import { Component, OnInit } from '@angular/core';
import { WasmGoService } from './services/wasm-go.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MermaidService } from './services/mermaid.service';
import { StateService } from './services/state.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent implements OnInit {

  protected mermaidContent: string = "";
  protected devfileYaml: string = ""

  constructor(
    protected sanitizer: DomSanitizer,
    private wasmGo: WasmGoService,
    private mermaid: MermaidService,
    private state: StateService,
  ) {}

  ngOnInit() {
    this.state.state.subscribe(async newContent => {
      if (newContent == null) {
        return;
      }

      this.devfileYaml = newContent.content;

      try {
        const result = this.wasmGo.getFlowChart();
        const svg = await this.mermaid.getMermaidAsSVG(result);
        this.mermaidContent = svg;  
      } catch {}
    });
  }

  async onButtonClick(content: string){
    const newContent = this.wasmGo.setDevfileContent(content);
    this.state.changeDevfileYaml(newContent);
  }

}
