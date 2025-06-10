import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * クラス名を結合し、Tailwindのユーティリティクラスの衝突を解消するユーティリティ関数
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

/**
 * HTMLエスケープを行うユーティリティ関数
 * XSS攻撃を防ぐために特殊文字をエスケープします
 */
export function escapeHtml(text: string | number): string {
  const escapeMap: { [key: string]: string } = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#039;'
  };
  return String(text).replace(/[&<>"']/g, (match) => escapeMap[match]);
}