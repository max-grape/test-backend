# Backend test project

## TODO:

- use shorten commit sha as docker image tag?  
  don't like crutches like that:
  ```
  - name: Set outputs
    id: vars
    run: echo "::set-output name=sha_short::$(git rev-parse --short HEAD)"
  - name: Check outputs
    run: echo ${{ steps.vars.outputs.sha_short }}
  ```
  ```
  - name: Set outputs
    run: echo "SHORT_SHA=`echo ${GITHUB_SHA} | cut -c1-8`" >> $GITHUB_ENV
  - name: Check outputs
    run: echo ${SHORT_SHA}
  ```
- reuse development workflow jobs in production workflow somehow
