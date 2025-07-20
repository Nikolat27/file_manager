// src/utils/toast.js
import { useToast } from "vue-toastification";

const toast = useToast();

export function showSuccess(message) {
    toast.success(message);
}

export function showError(message) {
    toast.error(message);
}

// Optionally export the raw toast if you want custom options
export { toast };
