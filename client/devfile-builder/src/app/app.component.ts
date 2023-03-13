import { Component } from '@angular/core';
import { WasmGoService } from './wasm-go.service';
import { DomSanitizer } from '@angular/platform-browser';
import { MermaidService } from './mermaid.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {

  protected mermaidContent: string = "";

  constructor(
    protected sanitizer: DomSanitizer,
    private wasmGo: WasmGoService,
    private mermaid: MermaidService,
  ) {}

  async onButtonClick(content: string){
    const result = this.wasmGo.getFlowChart(content);
    const svg = await this.mermaid.getMermaidAsSVG(result);
    this.mermaidContent = svg;
  }

}
