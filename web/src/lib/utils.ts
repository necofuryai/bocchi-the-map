import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * Utility function to combine class names and resolve Tailwind utility class conflicts
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * HTML escape map (defined outside function for performance improvement)
 */
const escapeMap: { [key: string]: string } = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  "'": '&#039;'
};

/**
 * Utility function to perform HTML escaping
 * Escapes special characters to prevent XSS attacks
 */
export function escapeHtml(text: string | number): string {
  return String(text).replace(/[&<>"']/g, (match) => escapeMap[match]);
}