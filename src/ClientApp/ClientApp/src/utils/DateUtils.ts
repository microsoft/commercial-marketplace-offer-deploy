// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.
export function toLocalDateTime(isoString: string): string {
    return new Date(isoString).toLocaleString();
  }
  