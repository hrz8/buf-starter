export const FALLBACK_RESPONSES = [
  'I\'m currently in maintenance mode. Please try again later.',
  'The AI service is temporarily unavailable. Please check back soon.',
  'I\'m taking a short break. Please try your request again in a few moments.',
  'Service is currently offline for updates. We\'ll be back shortly.',
  'I\'m not available right now, but I\'ll be back online soon!',
  'The system is undergoing maintenance. Thank you for your patience.',
  'AI capabilities are temporarily disabled. Please contact support if this persists.',
  'I\'m offline at the moment. Please try again later.',
  'Our AI assistant is currently unavailable. We apologize for the inconvenience.',
  'The service is temporarily down for improvements. Please check back in a few minutes.',
  'I\'m unable to process your request at the moment. Please try again shortly.',
  'The AI is taking a quick break. We\'ll be ready to help you soon!',
  'Service maintenance in progress. Thank you for your understanding.',
  'I\'m currently updating my systems. Please come back in a little while.',
  'The assistant is temporarily offline. Your patience is appreciated.',
  'We\'re making some improvements. Please try again in a few minutes.',
  'AI services are paused for maintenance. We\'ll be back shortly.',
  'I\'m not able to respond right now. Please check back soon.',
  'The system is temporarily unavailable. We\'re working to restore service.',
  'Our AI is currently offline. Please try your request again later.',
];

export function getRandomFallbackResponse(): string {
  const randomIndex = Math.floor(Math.random() * FALLBACK_RESPONSES.length);
  return FALLBACK_RESPONSES[randomIndex] ?? FALLBACK_RESPONSES[0] ?? 'Service unavailable.';
}
