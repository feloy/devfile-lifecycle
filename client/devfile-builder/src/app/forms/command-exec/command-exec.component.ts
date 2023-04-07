import { Component } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { StateService } from 'src/app/services/state.service';
import { WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-command-exec',
  templateUrl: './command-exec.component.html',
  styleUrls: ['./command-exec.component.css']
})
export class CommandExecComponent {
  form: FormGroup;
  containerList: string[] = [];

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {
    this.form = new FormGroup({
      name: new FormControl(""),
      component: new FormControl(),
      commandLine: new FormControl(""),
      workingDir: new FormControl(""),
      hotReloadCapable: new FormControl(false),
    });

    this.state.state.subscribe(async newContent => {
      const containers = newContent?.containers;
      if (containers == null) {
        return
      }
      this.containerList = containers.map(container => container.name);
    });
  }

  create() {
   console.log(this.form.value);
    const newDevfile = this.wasm.addExecCommand(this.form.value["name"], this.form.value);
    this.state.changeDevfileYaml(newDevfile);
  }
}
