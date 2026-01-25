import type { Project } from '~~/gen/altalune/v1/project_pb';

export const useProjectStore = defineStore('project', () => {
  const router = useRouter();
  const route = useRoute();

  const activeProjectId = useLocalStorage<string | null>('pId', null);
  const projects = ref<Project[]>([]);
  const pending = ref(false);
  const error = ref<Error | null>(null);
  const projectsLoaded = ref(false);

  const currentProject = computed(() => {
    if (!projects.value.length || !activeProjectId.value) {
      return null;
    }
    return projects.value.find(p => p.id === activeProjectId.value) || null;
  });

  const isProjectNotFound = computed(() => {
    return projectsLoaded.value
      && projects.value.length > 0
      && activeProjectId.value !== null
      && !projects.value.some(p => p.id === activeProjectId.value);
  });

  function setProjects(newProjects: Project[]) {
    projects.value = newProjects;
    projectsLoaded.value = true;

    if (!newProjects.length)
      return;

    const urlProjectId = route.query.pId as string | undefined;
    const storedProjectId = activeProjectId.value;

    const urlProjectValid = urlProjectId
      && newProjects.some(p => p.id === urlProjectId);
    const storedProjectValid = storedProjectId
      && newProjects.some(p => p.id === storedProjectId);

    if (urlProjectId && !urlProjectValid && storedProjectValid) {
      router.replace({
        query: {
          ...route.query,
          pId: storedProjectId,
        },
      });
      return;
    }

    if (urlProjectValid && urlProjectId !== storedProjectId) {
      activeProjectId.value = urlProjectId;
      return;
    }

    if (!storedProjectId && !urlProjectValid) {
      setActiveProject(newProjects[0]!.id);
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

  function addProject(project: Project) {
    projects.value = [project, ...projects.value];
    setActiveProject(project.id);
  }

  function removeProject(publicId: string) {
    projects.value = projects.value.filter(p => p.id !== publicId);
    if (projects.value.length > 0) {
      setActiveProject(projects.value[0]!.id);
    }
  }

  function setLoading(isLoading: boolean) {
    pending.value = isLoading;
  }

  function setError(err: Error | null) {
    error.value = err;
  }

  function selectProjectFromOverlay(projectId: string) {
    setActiveProject(projectId);
  }

  return {
    projects: readonly(projects),
    pending: readonly(pending),
    error: readonly(error),
    activeProjectId: readonly(activeProjectId),
    currentProject,
    isProjectNotFound,
    projectsLoaded: readonly(projectsLoaded),
    setProjects,
    setActiveProject,
    addProject,
    removeProject,
    setLoading,
    setError,
    selectProjectFromOverlay,
  };
});
