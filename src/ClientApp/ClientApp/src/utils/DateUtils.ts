export function toLocalDateTime(isoString: string): string {
    return new Date(isoString).toLocaleString();
  }
  