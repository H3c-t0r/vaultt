/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: BUSL-1.1
 */

import Component from '@glimmer/component';
import { pluralize } from 'ember-inflector';
import { inject as service } from '@ember/service';

/**
 * @module ListView
 * `ListView` components are used in conjunction with `ListItem` for rendering a list.
 *
 * @example
 * ```js
 * <ListView @items={{model}} @itemNoun="role" as |list|>
 *   {{#if list.empty}}
 *     <list.empty @title="No roles here" />
 *   {{else}}
 *     <div>
 *       {{list.item.id}}
 *     </div>
 *   {{/if}}
 * </ListView>
 * ```
 *
 * @param {array} [items=null] - An Ember array of items (objects) to render as a list. Because it's an Ember array it has properties like length an meta on it.
 * @param {string} [itemNoun=item] - A noun to use in the empty state of message and title.
 * @param {string} [message=null] - The message to display within the banner.
 * @param {boolean} [showPagination=false] - To show HDS pagination or not. If true, will show pagination even if only one item in the list.
 * @yields {object} Yields the current item in the loop.
 * @yields If there are no objects in items, then `empty` will be yielded - this is an instance of
 * the EmptyState component.
 * @yields If `item` or `empty` isn't present on the object, the component can still yield a block - this is
 * useful for showing states where there are items but there may be a filter applied that returns an
 * empty set.
 *
 */
export default class ListView extends Component {
  @service router;

  get itemNoun() {
    return this.args.itemNoun || 'item';
  }

  get emptyTitle() {
    const items = pluralize(this.itemNoun);
    return `No ${items} yet`;
  }

  get emptyMessage() {
    const items = pluralize(this.itemNoun);
    return `Your ${items} will be listed here. Add your first ${this.itemNoun} to get started.`;
  }

  // callback from HDS pagination to set the queryParams currentPage
  get paginationQueryParams() {
    return (page) => {
      return {
        page,
      };
    };
  }
}
