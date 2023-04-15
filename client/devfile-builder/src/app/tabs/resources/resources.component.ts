import { Component, OnInit } from '@angular/core';
import { StateService } from 'src/app/services/state.service';
import { ClusterResource } from 'src/app/services/wasm-go.service';

@Component({
  selector: 'app-resources',
  templateUrl: './resources.component.html',
  styleUrls: ['./resources.component.css']
})
export class ResourcesComponent implements OnInit {

  forceDisplayAdd: boolean = false;
  resources: ClusterResource[] | undefined = [];

  constructor(
    private state: StateService,
  ) {}

  ngOnInit() {
    const that = this;
    this.state.state.subscribe(async newContent => {
      console.log(newContent);
      that.resources = newContent?.resources;
      if (this.resources == null) {
        return
      }
      that.forceDisplayAdd = false;
    });
  }

  displayAddForm() {
    this.forceDisplayAdd = true;
  }

  undisplayAddForm() {
    this.forceDisplayAdd = false;
  }
}
