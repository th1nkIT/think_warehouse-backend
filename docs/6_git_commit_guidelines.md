#9. Git Commit Guidelines

copied from https://karma-runner.github.io/6.3/dev/git-commit-msg.html

## Format **commit message**

Use the following format:

    ```javascript
    <req_commit_type>(<req_commit_scope>): <req_commit_title>

    <optional_commit_description>

    <optional_commit_footer>
    ```

Where:

* `req_commit_type` (required) is one of:
  * `feat` = development of new features for users
  * `fix` = related to bugfixing source code in production
  * `style` = source code format, spaces, tabs, indents that have no effect
    to production deployment
  * `refactor` = refactor production code. example: rename variable, move logic to new function
  * `test` = addition of unit tests that have no effect on production
    deployment
  * `chore` = repair script, dockerfile, Makefile etc. Not affect
    production deployment
  * `docs` = changes to documentation, comments & README files

* `req_commit_scope` (required) i.e. scope of module / source code changes. Example:
  *config
  * logs
  * service
  * repository
  * postgres
    *echo
  * grpc

* `req_commit_title` (required) filled with a brief description of the changes that have been made

* `optional_commit_description` _(optional)_ more detailed description related
  source code changes. Can be filled with information why the change
  done (why & how), what is done (what).

* `optional_commit_footer` _(optional)_ link to JIRA ticket number, Gitlab &
  info `breaking changes API`.


### Example

* We apply changes related to ABC-123 issue, bugfix in endpoint login because people can
  login without sending device ID.

    ```javascript
    fix(login): add device ID validation to login endpoint

    wit.id/browse/ABC-123
    ```

* reformat source code from IDE (right click -> reformat code)

    ```javascript
    style(format): apply reformat code from IDE
    ```

* Issue PAY-123, add new features, payment via ForeignPay. Creating API
  client/service objects.

    ```javascript
    feat(asingpay): implement json http api client

    wit.id/browse/PAY-123
    ```

* Change the target Makefile, auto reformat code before run in local.

    ```javascript
    chore(Makefile): reformat code before run at local machine
    ```