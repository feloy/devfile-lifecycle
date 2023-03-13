import { Injectable } from '@angular/core';

declare const Go: any;
declare const getFlowChart: any;

@Injectable({
  providedIn: 'root'
})
export class WasmGoService {

  constructor() { 
    console.log("start wasm service");
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("./assets/devfile.wasm"), go.importObject).then((result) => {
        go.run(result.instance);                
    });
  }

  // getFlowChart calls the wasm module to get the lifecycle of the Devfile in mermaid chart format
  getFlowChart(devfile: string): string {
    return getFlowChart(devfile);
  }
}
