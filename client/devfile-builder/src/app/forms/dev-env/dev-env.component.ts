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
        return;
      }
      this.devEnvs().clear();
      for (const devEnv of newContent.devEnvs) {
        this.addDevEnv(devEnv);
      }
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

}
