<template>
  <v-container>
    <v-card>

      <v-card-title>
        Get User
      </v-card-title>
      <v-form
          ref="form"
          v-model="valid"
          @submit.prevent="submit"
          lazy-validation
      >
        <v-card-text>
          <v-text-field
              v-model="username"
              label="Username"
              :rules="[() => !!username || 'This field is required']"
              required
          ></v-text-field>

          <v-text-field
              v-model="firstname"
              label="Firstname"
              :rules="[() => !!firstname || 'This field is required']"
              required
          ></v-text-field>
          <v-text-field
              v-model="lastname"
              label="Lastname"

              :rules="[() => !!lastname || 'This field is required']"
              required
          ></v-text-field>

        </v-card-text>

        <v-card-actions>
          <v-btn
              :disabled="!valid"
              class="mr-4"
              @click="submit"
          >
            Login
          </v-btn>

        </v-card-actions>

      </v-form>
    </v-card>
  </v-container>

</template>

<script>
export default {
  username: "FakeLogin",
  data: () => ({
    valid: false,
    username: '',
    firstname: '',
    lastname: '',
  }),

  methods: {
    submit () {
      this.$refs.form.validate()
      this.$store.dispatch('login', {
        username: this.username,
        firstname: this.firstname,
        lastname: this.lastname
      })
    },

  },
}
</script>

<style scoped>

</style>