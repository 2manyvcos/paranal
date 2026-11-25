export class HTTPError extends Error {
  readonly status: number;

  constructor(msg: string, status: number) {
    super(msg);
    Object.setPrototypeOf(this, HTTPError.prototype);
    this.status = status;
  }
}

export const HTTP_NOT_FOUND = 404;
