import { Request, Response, Router } from 'express';
import * as jwt from 'express-jwt'
import { User } from '../domain/user/User';
import { UserManagement } from '../domain/user/UserManagement';
import { Controller } from './Controller';

export class UserController extends Controller {

  public basePath: string = '/users'

  public constructor(
    public userManagement: UserManagement
  ) {
    super()
  }

  public getRoutes(checkJwt: jwt.RequestHandler): Router {
    const router = Router()
    router.post('/', checkJwt, this.handleConnection.bind(this))
    router.get('/:email', this.getUser.bind(this))
    return router
  }

  public getUser(req: Request, res: Response) {
    this.userManagement
      .findOne(req.params.email)
      .then(res.json)
      .catch(err => {
        res.json(err)
      })
  }

  public handleConnection(req: Request, res: Response) {
    this.userManagement
      .create(new User(req.user.email))
      .then(res.json)
      .catch(res.json)
  }

}