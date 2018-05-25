# Creative Cloud Daemon (ccd)

Creative Cloud Daemon (CCD) project allows Linux users to
leverage Creative Cloud storage services. CCD is a
storage client which enables users to upload files to and
download them from the Adobe Creative Cloud. CCD is also
a daemon, that when running locally, allows users to
interact with the Creative Cloud backend and integrate
workflows with Adobe Cloud Services.

The project is under development. Please feel free to
make contributions and join us.


### Testing

As so to run the unit tests you can use the make task test as follows:

```
make test
```

Furthermore, integration tests need to communicate with an actual backend and
required valid credentials for the environment in which the tests are run:

```
make integration-test
```

The project does leverage "build tags" to distinguish the unit tests from the
integration tests. You are required to mark all tests with "all" tag and the
integration tests with "all" and "integration-test". 

