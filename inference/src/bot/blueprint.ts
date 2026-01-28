import type { BotRepository } from '../repositories/db/bot.js';
import type { ProjectRepository } from '../repositories/db/project.js';
import type { BlueprintData } from './types.js';

export class Blueprint {
  private static projectRepo: ProjectRepository;
  private static botRepo: BotRepository;

  private constructor() {}

  public static setup(projectRepo: ProjectRepository, botRepo: BotRepository): void {
    this.projectRepo = projectRepo;
    this.botRepo = botRepo;
  }

  public static async load(projectPublicId: string): Promise<BlueprintData | null> {
    const projectInfo = await this.projectRepo.getProjectInfo(projectPublicId);
    if (!projectInfo) {
      console.error(`project not found: ${projectPublicId}`);
      return null;
    }

    const [modulesConfig, nodes] = await Promise.all([
      this.botRepo.getModulesConfig(projectInfo.id),
      this.botRepo.getNodes(projectInfo.id),
    ]);

    return {
      projectName: projectInfo.name,
      modules: modulesConfig ?? {},
      nodes,
    };
  }
}
