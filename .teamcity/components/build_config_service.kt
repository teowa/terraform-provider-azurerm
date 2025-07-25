import jetbrains.buildServer.configs.kotlin.*

class serviceDetails(name: String, displayName: String, environment: String, vcsRootId : String) {
    val packageName = name
    val displayName = displayName
    val environment = environment
    val vcsRootId = vcsRootId

    fun buildConfiguration(providerName : String, nightlyTestsEnabled: Boolean, startHour: Int, parallelism: Int, daysOfWeek: String, daysOfMonth: String, timeout: Int, disableTriggers: Boolean) : BuildType {
        return BuildType {
            // TC needs a consistent ID for dynamically generated packages
            id(uniqueID(providerName))

            name = "%s - Acceptance Tests".format(displayName)

            vcs {
                root(rootId = AbsoluteId(vcsRootId))
                cleanCheckout = true
            }

            steps {
                ConfigureGoEnv()
                DownloadTerraformBinary()
                RunAcceptanceTests(packageName)
            }

            failureConditions {
                errorMessage = true
                executionTimeoutMin = 60 * timeout
            }

            features {
                Golang()
            }

            params {
                TerraformAcceptanceTestParameters(parallelism, "TestAcc", timeout)
                TerraformAcceptanceTestsFlag()
                TerraformCoreBinaryTesting()
                TerraformShouldPanicForSchemaErrors()
                ReadOnlySettings()
                WorkingDirectory(packageName)
            }

            triggers {
                RunNightly(nightlyTestsEnabled, startHour, daysOfWeek, daysOfMonth, disableTriggers)
            }
        }
    }

    fun uniqueID(provider : String) : String {
        return "%s_SERVICE_%s_%s".format(provider.uppercase(), environment.uppercase(), packageName.uppercase())
    }
}
