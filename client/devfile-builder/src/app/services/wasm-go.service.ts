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
  commands: Command[];
  containers: Container[];
  devEnvs: DevEnv[];
};

export type Metadata = {
  name: string | null;
  version: string | null;
  displayName: string | null;
  description: string | null;
  tags: string | null;
  architectures: string | null;
  icon: string | null;
  globalMemoryLimit: string | null;
  projectType: string | null;
  language: string | null;
  website: string | null;
  provider: string | null;
  supportUrl: string | null;
};

export type Command = {
  name: string;
  type: "exec" | "apply" | "composite";
  exec: ExecCommand | undefined;
  apply: ApplyCommand | undefined;
  composite: CompositeCommand | undefined;
};

export type ExecCommand = {
  component: string;
  commandLine: string;
  workingDir: string;
  hotReloadCapable: boolean;
};

export type ApplyCommand = {
  component: string;
};

export type CompositeCommand = {
  commands: string;
  parallel: boolean;
};

export type Container = {
  name: string;
  image: string;
  command: string[];
  args: string[];
};

export type DevEnv = {
  name: string;
  image: string;
  command: string[];
  args: string[];
  userCommands: UserCommand[];
}

export type Group = '' | 'build' | 'test'| 'run'  | 'debug' | 'deploy';

export type UserCommand = {
  name: string;
  group: Group;
  default: boolean;
  commandLine: string;
  hotReloadCapable: boolean;
  workInSourceDir: boolean;
  workingDir: string;
};

declare const addContainer: (name: string, image: string, command: string[], args: string[]) => Result;
declare const addUserCommand: (component: string, name: string, commandLine: string) => Result;
declare const getFlowChart: () => ChartResult;
declare const setDevfileContent: (devfile: string) => Result;
declare const setMetadata: (metadata: Metadata) => Result;
declare const updateContainer: (name: string, image: string, command: string[], args: string[], userCommands: UserCommand[]) => Result;

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

  addContainer(container: Container): ResultValue {
    console.log("container", container);
    const result = addContainer(
      container.name,
      container.image,
      container.command,
      container.args,
    );
    return result.value;
  }

  addUserCommand(component: string, name: string, commandLine: string): ResultValue {
    const result = addUserCommand(component, name, commandLine);
    return result.value;
  }

  updateContainer(devEnv: DevEnv): ResultValue {
    const result = updateContainer(
      devEnv.name,
      devEnv.image,
      devEnv.command,
      devEnv.args,
      devEnv.userCommands,
    );
    if (result.err != "") {
      console.log(result.err);
    }
    return result.value;  
  }

  // getFlowChart calls the wasm module to get the lifecycle of the Devfile in mermaid chart format
  getFlowChart(): string {
    const result = getFlowChart();
    return result.value;
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

}
