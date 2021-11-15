/**
 * @module <%= classifiedModuleName %>
 * <%= classifiedModuleName %> components are used to...
 *
 * @example
 * ```js
 * <<%= classifiedModuleName %> @requiredParam={requiredParam} @optionalParam={optionalParam} @param1={{param1}}/>
 * ```
 * @param {object} requiredParam - requiredParam is...
 * @param {string} [optionalParam] - optionalParam is...
 * @param {string} [param1=defaultValue] - param1 is...
 */


import Component from '@glimmer/component';
import { setComponentTemplate } from '@ember/component';
<%= importTemplate %>

class <%= classifiedModuleName %> extends Component {<%= contents %>  
}

export default setComponentTemplate(<%= hasLayout %>, <%= classifiedModuleName %>)