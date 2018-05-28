import { Controller, Get, Route } from "tsoa"
import Logger from '../Logger'

@Route("_health")
export class HealthController extends Controller {
  
  constructor() {
    super()
  }

  @Get("")
  public handleHealthCheck(): Promise<any> {
    this.setStatus(200)
    return Promise.resolve({status: 'ok'})
  }

}
