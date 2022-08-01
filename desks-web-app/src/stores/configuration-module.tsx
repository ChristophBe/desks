import {ActionContext, ActionTree, Module, MutationTree} from "vuex";
import {RootState} from "@/stores/store";
import FrontendConfiguration from "@/dtos/FrontendConfigurationDto";
import {getData} from "@/utils/request-utils";

export interface ConfigurationState {
    loading: boolean
    configuration: FrontendConfiguration | null
}

const state: ConfigurationState = {
    loading: false,
    configuration: null
}

const mutations: MutationTree<ConfigurationState> = {
    loading(state: ConfigurationState) {
        state.loading = true
    },
    setConfiguration(state: ConfigurationState, configuration: FrontendConfiguration) {
        state.loading = false
        state.configuration = configuration
    }
}

const actions: ActionTree<ConfigurationState, RootState> = {

    async fetchConfiguration({commit}: ActionContext<ConfigurationState, RootState>) {

        const response = await getData(`/api/v1.0/configuration`)
        if (response.status >= 400) {
            console.log("Failed to fetch configuration", response.status)
            return
        }
        try {
            const configuration = await response.json()
            commit('setConfiguration', configuration)
        } catch (e) {
            console.log("Failed to fetch configuration", e)
        }

    },
}
export const configurationModule: Module<ConfigurationState, RootState> = {
    namespaced: true,
    state,
    mutations,
    actions
}

