import { Injectable } from '@angular/core';

import { BehaviorSubject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class StateService {

  private _state = new BehaviorSubject<string>("");
  public state = this._state.asObservable(); 

  changeDevfileYaml(newYaml: string) {
    this._state.next(newYaml);
  }
}
