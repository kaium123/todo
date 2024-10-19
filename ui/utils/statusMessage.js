export function setStatusMessage(vueInstance, message) {
    vueInstance.statusMessage = message;
    setTimeout(() => {
      vueInstance.statusMessage = '';
    }, 5000);
  }
  