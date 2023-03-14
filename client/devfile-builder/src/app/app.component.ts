import { Component, OnInit } from '@angular/core';
import { WasmGoService } from './wasm-go.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MermaidService } from './mermaid.service';
import { StateService } from './state.service';

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
      this.state.state.subscribe(async newYaml => {
        if (newYaml == "") {
          return;
        }
        console.log(newYaml);
        const result = this.wasmGo.getFlowChart();
        const svg = await this.mermaid.getMermaidAsSVG(result);
        this.mermaidContent = svg;
  
        this.devfileYaml = newYaml;
      });
    }

  async onButtonClick(content: string){
    this.wasmGo.setDevfileContent(content);
    this.state.changeDevfileYaml(content);
  }

}
