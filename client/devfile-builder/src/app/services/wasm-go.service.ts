import { Injectable } from '@angular/core';

declare const Go: any;

type ChartResult = {
  err: string;
  value: any;
};

type Result = {
  err: string;
  value: ResultValue;
};

export type ResultValue = {
  content: string;
  metadata: Metadata;
};

type Metadata = {
  name: string | null;
  displayName: string | null;
  description: string | null;
};

declare const getFlowChart: () => ChartResult;
declare const setDevfileContent: (devfile: string) => Result;
declare const setMetadata: (metadata: Metadata) => Result;

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
  setDevfileContent(devfile: string): ResultValue {
    const result = setDevfileContent(devfile);
    return result.value;
  }

  setMetadata(metadata: Metadata): ResultValue {
    const result = setMetadata(metadata);
    return result.value;
  }

  // getFlowChart calls the wasm module to get the lifecycle of the Devfile in mermaid chart format
  getFlowChart(): string {
    const result = getFlowChart();
    return result.value;
  }
}
