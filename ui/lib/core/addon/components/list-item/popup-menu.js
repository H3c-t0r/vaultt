/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

import Component from '@ember/component';
import layout from '../../templates/components/list-item/popup-menu';

export default Component.extend({
  layout,
  tagName: '',
  item: null,
  hasMenu: null,
});
