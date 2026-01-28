import type {
  BlueprintData,
  ChatbotNode,
  ModuleConfigMap,
  ModuleName,
  PlainMessage,
} from './types.js';
import { Blueprint } from './blueprint.js';

export class Bot {
  public modules = new Map<ModuleName, ModuleConfigMap[ModuleName]>();
  public nodes = new Map<string, PlainMessage<ChatbotNode>>();

  public constructor(public readonly projectId: string) {}

  public async load(): Promise<void> {
    const data = await Blueprint.load(this.projectId);
    if (!data) {
      console.error(`failed to load blueprint for project: ${this.projectId}`);
      return;
    }

    console.info(`loading bot for '${data.projectName}'`);
    this.loadBluePrint(data);
  }

  public loadBluePrint(data: BlueprintData): void {
    const { modules, nodes } = data;
    this.modules.clear();
    this.nodes.clear();

    (Object.keys(modules) as ModuleName[])
      .filter(moduleId => modules[moduleId]?.enabled)
      .forEach((moduleId) => {
        const moduleConfig = modules[moduleId];
        if (moduleConfig) {
          this.loadModule(moduleId, moduleConfig);
          console.info(`module '${moduleId}' loaded`);
        }
      });

    const nodeEntries = nodes instanceof Map
      ? nodes.entries()
      : Object.entries(nodes);
    for (const [nodeId, node] of nodeEntries) {
      this.nodes.set(nodeId, node);
      console.info(`node '${node.name}' loaded`);
    }
  }

  public loadModule<K extends ModuleName>(moduleId: K, options: ModuleConfigMap[K]): void {
    this.modules.set(moduleId, options);
  }

  // public getModule<K extends ModuleName>(moduleId: K): ModuleConfigMap[K] | undefined {
  //   return this.modules.get(moduleId) as ModuleConfigMap[K] | undefined;
  // }

  // public getNode(nodeId: string): PlainMessage<ChatbotNode> | undefined {
  //   return this.nodes.get(nodeId);
  // }

  // public isModuleEnabled(moduleId: ModuleName): boolean {
  //   const module = this.modules.get(moduleId);
  //   return module?.enabled ?? false;
  // }
}
