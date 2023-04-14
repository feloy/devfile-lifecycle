import { Component } from '@angular/core';
import { FormControl, FormGroup } from '@angular/forms';
import { StateService } from 'src/app/services/state.service';
import { WasmGoService } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-command-image',
  templateUrl: './command-image.component.html',
  styleUrls: ['./command-image.component.css']
})
export class CommandImageComponent {
  form: FormGroup;
  imageList: string[] = [];

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {
    this.form = new FormGroup({
      name: new FormControl(),
      component: new FormControl(),
    });

    this.state.state.subscribe(async newContent => {
      const images = newContent?.images;
      if (images == null) {
        return
      }
      this.imageList = images.map(image => image.name);
    });
  }

  create() {
    console.log(this.form.value);
    const newDevfile = this.wasm.addApplyCommand(this.form.value["name"], this.form.value);
    this.state.changeDevfileYaml(newDevfile);
  }
}