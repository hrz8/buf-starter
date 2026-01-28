export class Bot {
  public modules = new Map<string, any>();
  public nodes = new Map<string, any>();

  constructor(public readonly projectId: string) {
  }

  public async load(): Promise<void> {
    console.info(`loading bot for: ${this.projectId}`);

    const dataFromDB = {
      modules: { llm: { enabled: true } },
      nodes: {},
    };
    this.loadBluePrint(dataFromDB);
  }

  public loadBluePrint(data: any): void {
    const {
      modules,
      // eslint-disable-next-line unused-imports/no-unused-vars
      nodes,
    } = data;
    this.modules.clear();
    this.nodes.clear();

    Object.keys(modules).filter(moduleId => modules[moduleId].enabled).forEach((moduleId) => {
      this.loadModule(moduleId);
      console.info(`Module '${moduleId}' loaded`);
    });
  }

  public loadModule(moduleId: string): void {
    this.modules.set(moduleId, {});
  }
}
