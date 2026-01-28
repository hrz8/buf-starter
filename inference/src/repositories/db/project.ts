import type { Knex } from 'knex';

interface ProjectRow {
  id: string;
  public_id: string;
  name: string;
}

export interface ProjectInfo {
  id: string;
  name: string;
}

export class ProjectRepository {
  private readonly db: Knex;

  public constructor(db: Knex) {
    this.db = db;
  }

  public async getProjectInfo(publicId: string): Promise<ProjectInfo | null> {
    const project = await this.db<ProjectRow>('altalune_projects')
      .select('id', 'public_id', 'name')
      .where('public_id', publicId)
      .first();

    if (!project) {
      return null;
    }

    return {
      id: project.id,
      name: project.name,
    };
  }
}
