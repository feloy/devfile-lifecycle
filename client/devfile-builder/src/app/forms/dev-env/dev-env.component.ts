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
      console.log(newContent?.devEnvs)
      
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
        this.addDevEnv(devEnv);
      }
      this.showCreate = false;
    });
  }

  newDevEnv(devEnv: DevEnv): FormGroup {
    return new FormGroup({
      name: new FormControl(devEnv.name),
      image: new FormControl(devEnv.image),
    })
  }

  addDevEnv(devEnv: DevEnv) {
    this.devEnvs().push(this.newDevEnv(devEnv))
  }

  devEnvs(): FormArray {
    return this.form.get('devEnvs') as FormArray;
  }

  save(i: number) {
    console.log(i);
    const devEnvToSave = this.form.value['devEnvs'][i];
    console.log(devEnvToSave);
  }

  createNew() {
    console.log(this.newName.value);
    console.log(this.newImage.value);
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
