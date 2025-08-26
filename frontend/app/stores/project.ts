import type { Project } from '~~/gen/altalune/v1/project_pb';

export const useProjectStore = defineStore('project', () => {
  const router = useRouter();
  const route = useRoute();

  const activeProjectId = useLocalStorage<string | null>('pId', null);
  const projects = ref<Project[]>([]);
  const pending = ref(false);
  const error = ref<Error | null>(null);

  const currentProject = computed(() => {
    if (!projects.value.length || !activeProjectId.value) {
      return null;
    }
    return projects.value.find((p) => p.id === activeProjectId.value) || null;
  });

  function setProjects(newProjects: Project[]) {
    projects.value = newProjects;

    if (newProjects.length && !currentProject.value) {
      const urlProjectId = route.query.pId as string;
      const targetProject = urlProjectId
        ? newProjects.find((p) => p.id === urlProjectId)
        : newProjects[0];

      if (targetProject) {
        setActiveProject(targetProject.id);
      }
    }
  }

  function setActiveProject(pId: string) {
    activeProjectId.value = pId;

    if (pId && route.query.pId !== pId) {
      router.replace({
        query: {
          ...route.query,
          pId,
        },
      });
    }
  }

  function setLoading(isLoading: boolean) {
    pending.value = isLoading;
  }

  function setError(err: Error | null) {
    error.value = err;
  }

  return {
    projects: readonly(projects),
    pending: readonly(pending),
    error: readonly(error),
    activeProjectId: readonly(activeProjectId),
    currentProject,
    setProjects,
    setActiveProject,
    setLoading,
    setError,
  };
});
