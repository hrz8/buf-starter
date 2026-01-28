export class AppConfig {
  private readonly connectionString
    = 'postgres://postgres:toor@localhost:5432/altalune?sslmode=disable';

  private readonly defaultProjectId = 'lb5pzkgrnbanlw';

  public constructor() {}

  public getDefaultProjectId(): string {
    return this.defaultProjectId;
  }

  public getConnectionString(): string {
    return this.connectionString;
  }
}
