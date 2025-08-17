import { ErrorDetailSchema } from '~~/gen/altalune/v1/common_pb';
import { ConnectError } from '@connectrpc/connect';

export function useErrorMessage() {
  const { t } = useI18n();

  function parseError(err: unknown): string {
    if (!(err instanceof ConnectError)) return t('errorCodes.6999');

    const detail = err.findDetails(ErrorDetailSchema)[0];
    if (!detail || !detail.code) return err.rawMessage;

    const code = detail.code.toString();
    const meta = detail.meta ?? {};

    return t(`errorCodes.${code}`, meta);
  }

  return {
    parseError,
  };
}
