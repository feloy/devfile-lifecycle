import { Injectable } from '@angular/core';

declare const Go: any;
declare const getFlowChart: any;
declare const setDevfileContent: any;

@Injectable({
  providedIn: 'root'
})
// WasmGoService uses the wasm module. 
// The module manages a single instance of a Devfile
export class WasmGoService {

  constructor() { 
    console.log("start wasm service");
    const go = new Go();
    WebAssembly.instantiateStreaming(fetch("./assets/devfile.wasm"), go.importObject).then((result) => {
        go.run(result.instance);                
    });
  }

  // setDevfileContent calls the wasm module to reset the content of the Devfile
  setDevfileContent(devfile: string) {
    setDevfileContent(devfile);
  }

  // getFlowChart calls the wasm module to get the lifecycle of the Devfile in mermaid chart format
  getFlowChart(): string {
    return getFlowChart();
  }
}
