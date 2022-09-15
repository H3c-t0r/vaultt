// jscodeshift can take a parser, like "babel", "babylon", "flow", "ts", or "tsx"
// Read more: https://github.com/facebook/jscodeshift#parser
// export const parser = 'babylon';

export default function transformer({ source }, api) {
  const j = api.jscodeshift;
  const filterForStore = (path) => {
    return j(path.value).find(j.MemberExpression, {
      object: {
        type: 'ThisExpression',
      },
      property: {
        name: 'store',
      },
    }).length;
  };
  let didInjectStore = false;

  // find class bodies and filter down to ones that access this.store
  const classesAccessingStore = j(source).find(j.ClassBody).filter(filterForStore);

  if (classesAccessingStore.length) {
    // filter down to class bodies where service is not injected
    const missingService = classesAccessingStore.filter((path) => {
      return !j(path.value)
        .find(j.ClassProperty, {
          key: {
            name: 'store',
          },
        })
        .filter((path) => {
          // ensure store property belongs to service decorator
          return path.value.decorators.find((path) => path.expression.name === 'service');
        }).length;
    });

    if (missingService.length) {
      // inject store service
      const storeService = j.classProperty(j.identifier('@service store'), null);
      // adding a decorator this way will force store down to a new line and then add a new line
      // leaving in just in case it's needed
      // storeService.decorators = [j.decorator(j.identifier('service'))];

      source = missingService
        .forEach((path) => {
          path.value.body.unshift(storeService);
        })
        .toSource();

      didInjectStore = true;
    }
  }

  // find object expressions and filter down to ones that access this.store (Ember classic)
  const objectsAccessingStore = j(source).find(j.ObjectExpression).filter(filterForStore);

  if (objectsAccessingStore.length) {
    // filter down to objects where service is not injected
    const missingService = objectsAccessingStore.filter((path) => {
      return !j(path.value).find(j.ObjectProperty, {
        key: {
          name: 'store',
        },
        value: {
          callee: {
            name: 'service',
          },
        },
      }).length;
    });

    if (missingService.length) {
      // inject store service
      const storeService = j.objectProperty(
        j.identifier('store'),
        j.callExpression(j.identifier('service'), [])
      );

      source = missingService
        .forEach((path) => {
          path.value.properties.unshift(storeService);
        })
        .toSource();

      didInjectStore = true;
    }
  }

  // if store was injected here check if inject has been imported
  if (didInjectStore) {
    const needsImport = !j(source).find(j.ImportSpecifier, {
      imported: {
        name: 'inject',
      },
    }).length;

    if (needsImport) {
      const injectionImport = j.importDeclaration(
        [j.importSpecifier(j.identifier('inject'), j.identifier('service'))],
        j.literal('@ember/service')
      );

      const imports = j(source).find(j.ImportDeclaration);
      source = imports
        .at(imports.length - 1)
        .insertAfter(injectionImport)
        .toSource({ reuseWhitespace: false });
    }
  }

  return source;
}
