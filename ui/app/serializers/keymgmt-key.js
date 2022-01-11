import ApplicationSerializer from './application';

export default class KeymgmtKeySerializer extends ApplicationSerializer {
  normalizeItems(payload) {
    let normalized = super.normalizeItems(payload);
    // Transform versions from object with number keys to array with key ids
    if (normalized.versions) {
      let lastRotated;
      let created;
      let versions = [];
      Object.keys(normalized.versions).forEach((key, i, arr) => {
        versions.push({
          id: parseInt(key, 10),
          ...normalized.versions[key],
        });
        if (i === 0) {
          created = normalized.versions[key].creation_time;
        } else if (arr.length - 1 === i) {
          // Set lastRotated to the last key
          lastRotated = normalized.versions[key].creation_time;
        }
      });
      normalized.versions = versions;
      return { ...normalized, last_rotated: lastRotated, created };
    } else if (Array.isArray(normalized)) {
      return normalized.map((key) => ({
        id: key.id,
        name: key.id,
        backend: payload.backend,
      }));
    }
    return normalized;
  }
}
