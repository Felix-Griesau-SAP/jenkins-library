import groovy.transform.Field
import com.sap.piper.Utils
import com.sap.piper.ConfigurationHelper

import static com.sap.piper.Prerequisites.checkScript

@Field String STEP_NAME = getClass().getName()
@Field Set GENERAL_CONFIG_KEYS = [
    /** Creates a SAP BTP ABAP Environment system via the cloud foundry command line interface */
    'abapEnvironmentCreateSystem',
    /** Deletes a SAP BTP ABAP Environment system via the cloud foundry command line interface */
    'cloudFoundryDeleteService',
    /** If set to true, a confirmation is required to delete the system */
    'confirmDeletion',
    'debug', // If set to true, the system is never deleted
    'testBuild', // Parameter for test execution mode, if true stage will be skipped
    'integrationTestOption' // Integration test option
]
@Field Set STAGE_STEP_KEYS = GENERAL_CONFIG_KEYS
@Field Set STEP_CONFIG_KEYS = STAGE_STEP_KEYS
@Field Set PARAMETER_KEYS = STEP_CONFIG_KEYS
/**
 * This stage creates a system for Integration Tests. The (custom) tests themselves can be added via a stage extension.
 */
void call(Map parameters = [:]) {
    def script = checkScript(this, parameters) ?: this
    def stageName = parameters.stageName?:env.STAGE_NAME

    Map config = ConfigurationHelper.newInstance(this)
        .loadStepDefaults()
        .mixinGeneralConfig(script.commonPipelineEnvironment, GENERAL_CONFIG_KEYS)
        .mixinStageConfig(script.commonPipelineEnvironment, stageName, STEP_CONFIG_KEYS)
        .mixin(parameters, PARAMETER_KEYS)
        .addIfEmpty('confirmDeletion', false)
        .addIfEmpty('debug', false)
        .addIfEmpty('testBuild', false)
        .addIfEmpty('integrationTestOption', 'addOnDeployment')
        .use()

    if (config.testBuild) {
        echo "Stage 'Integration Tests' skipped as parameter 'testBuild' is active"
        return null;
    }
    piperStageWrapper (script: script, stageName: stageName, stashContent: [], stageLocking: false) {
        if (config.integrationTestOption == 'systemProvisioning') {
            try {
                abapEnvironmentCreateSystem(script: parameters.script, includeAddon: true)
                cloudFoundryCreateServiceKey(script: parameters.script)
                abapEnvironmentBuild(script: parameters.script, phase: 'GENERATION', downloadAllResultFiles: true, useFieldsOfAddonDescriptor: '[{"use":"Name","renameTo":"SWC"}]')
            } catch (Exception e) {
                echo "Deployment test of add-on product failed."
                throw e
            } finally {
                if (!config.debug) {
                    cloudFoundryDeleteService script: parameters.script
                }
            }
        } else if (config.integrationTestOption == 'addOnDeployment') {
            try {
                abapLandscapePortalUpdateAddOnProduct(script: parameters.script)
                // abapEnvironmentBuild(script: parameters.script, phase: 'GENERATION', downloadAllResultFiles: true, useFieldsOfAddonDescriptor: '[{"use":"Name","renameTo":"SWC"}]')
            } catch (Exception e) {
                echo "Deployment test of add-on product failed."
                throw e
            }
        } else {
            e = new Error('Unsupoorted integration test option.')
            throw e
        }

        if (config.confirmDeletion) {
            input message: "Deployment test has been executed. Once you proceed, the test system will be deleted."
        }
    }
}
