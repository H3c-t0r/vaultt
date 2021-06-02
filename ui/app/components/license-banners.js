/**
 * @module LicenseBanners
 * LicenseBanners components are used to display Vault-specific license expiry messages
 *
 * @example
 * ```js
 * <LicenseBanners @expiry={expiryDate} />
 * ```
 * @param {string} expiry - RFC3339 date timestamp
 */

import Component from '@glimmer/component';
import isAfter from 'date-fns/isAfter';
import differenceInDays from 'date-fns/differenceInDays';

export default class LicenseBanners extends Component {
  get licenseExpired() {
    if (!this.args.expiry) return false;
    const now = new Date();
    return isAfter(now, new Date(this.args.expiry));
  }

  get licenseExpiringInDays() {
    if (!this.args.expiry) return -1;
    const now = new Date();
    return differenceInDays(new Date(this.args.expiry), now);
  }
}
