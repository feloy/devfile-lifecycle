import { Component, OnInit } from '@angular/core';
import { FormControl } from '@angular/forms';
import { StateService } from 'src/app/services/state.service';
import { WasmGoService } from 'src/app/services/wasm-go.service';


@Component({
  selector: 'app-metadata',
  templateUrl: './metadata.component.html',
  styleUrls: ['./metadata.component.css']
})
export class MetadataComponent implements OnInit {
  name = new FormControl('');
  displayName = new FormControl('');
  description = new FormControl('');

  constructor(
    private wasm: WasmGoService,
    private state: StateService,
  ) {}

  ngOnInit() {
    this.state.state.subscribe(async newContent => {
      const metadata = newContent?.metadata;
      if (metadata == null) {
        return
      }
      this.name.setValue(metadata.name);
      this.displayName.setValue(metadata.displayName);
      this.description.setValue(metadata.description);
    });
  }

  onSave() {
    const newDevfile = this.wasm.setMetadata({
      name: this.name.value,
      displayName: this.displayName.value,
      description: this.description.value,
    })
    this.state.changeDevfileYaml(newDevfile);
  }
}
