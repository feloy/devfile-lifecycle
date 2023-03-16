import { Component, OnInit } from '@angular/core';
import { FormArray, FormControl, FormGroup } from '@angular/forms';

import { Observable } from 'rxjs';

import { StateService } from 'src/app/services/state.service';
import { DevEnv, WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-dev-env',
  templateUrl: './dev-env.component.html',
  styleUrls: ['./dev-env.component.css']
})
export class DevEnvComponent implements OnInit {

  form: FormGroup;

  showCreate: boolean = false;
  newName = new FormControl('');
  newImage = new FormControl('');

   constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {
    this.form = new FormGroup({
      devEnvs: new FormArray([])
    })
  }

  ngOnInit() {
    this.state.state.subscribe(async newContent => {
      
      if (!newContent) {
        this.showCreate = true;
        return;
      }

      if (!newContent.devEnvs.length) {
        this.showCreate = true;
        return;
      }

      this.devEnvs().clear();
      for (const devEnv of newContent.devEnvs) {
        const devEnv_i = this.addDevEnv(devEnv);

        for (const command of devEnv.command) {
          this.addCommand(devEnv_i, command);
        }

        for (const arg of devEnv.args) {
          this.addArg(devEnv_i, arg);
        }
      }
      this.showCreate = false;
    });
  }

  newDevEnv(devEnv: DevEnv): FormGroup {
    return new FormGroup({
      name: new FormControl(devEnv.name),
      image: new FormControl(devEnv.image),
      command: new FormArray([]),
      args: new FormArray([]),
    })
  }

  addDevEnv(devEnv: DevEnv): number {
    this.devEnvs().push(this.newDevEnv(devEnv));
    return this.devEnvs().length-1;
  }

  devEnvs(): FormArray {
    return this.form.get('devEnvs') as FormArray;
  }

  addCommand(devEnv_i: number, command: string) {
    this.commands(devEnv_i).push(new FormControl(command));
  }

  commands(devEnv_i: number): FormArray {
    return this.devEnvs().controls[devEnv_i].get('command') as FormArray;
  }

  addArg(devEnv_i: number, arg: string) {
    this.args(devEnv_i).push(new FormControl(arg));
  }

  args(devEnv_i: number): FormArray {
    return this.devEnvs().controls[devEnv_i].get('args') as FormArray;
  }

  save(i: number) {
    const devEnvToSave = this.form.value['devEnvs'][i];
    console.log(devEnvToSave);
  }

  createNew() {
    if (this.newName.value == null || this.newImage.value == null) {
      // TODO should not happen with form validation
      return;
    }
    const newDevfile = this.wasm.addContainer(this.newName.value, this.newImage.value);
    this.state.changeDevfileYaml(newDevfile);
    
    this.resetNew();
  }

  resetNew() {
    this.newName.setValue("");
    this.newImage.setValue("");
    this.showCreate = false;
  }

  createCancel() {
    this.showCreate = false;
  }
}
