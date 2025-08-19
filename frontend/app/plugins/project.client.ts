import { useProjectStore } from '~/stores/project';

export default defineNuxtPlugin(() => {
  const router = useRouter();
  const route = useRoute();

  const projectStore = useProjectStore();

  router.beforeEach((to) => {
    const currentPId = (route.query.pId as string) ?? projectStore.activeProjectId;
    const toHasPId = 'pId' in to.query;

    if (currentPId && !toHasPId) {
      return {
        path: to.path,
        query: {
          ...to.query,
          pId: currentPId,
        },
      };
    }
  });
});
