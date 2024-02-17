<template>
  <v-container fluid class="fill-height">
    <!-- UPLOAD PAGE -->
    <v-row v-if="state === 'upload'" class="fill-height">
      <v-col class="text-center" cols="12">
        <v-btn
          block
          class="h-100"
          size="x-large"
          :color="btnColor"
          :loading="isLoading"
          :disabled="isLoading"
          @click="chooseFile"
        >
          Foto auswählen
        </v-btn>
        <v-file-input
          class="d-none"
          id="fileUpload"
          v-model="myFile"
          @change="showSelectedFile"
          accept="image/*,.heic,.heif"
        ></v-file-input>
      </v-col>
    </v-row>

    <!-- PREVIEW PAGE -->
    <v-row
      v-else-if="state === 'preview'"
      class="fill-height"
      align-content="space-between"
    >
      <v-col cols="12">
        <v-img :src="imageSrc" class="preview"></v-img>
      </v-col>
      <v-col cols="12">
        <v-btn
          size="x-large"
          block
          color="red"
          :disabled="isLoading"
          @click="state = 'upload'"
        >
          Zurück
        </v-btn>
        <v-btn
          class="mt-4"
          size="x-large"
          block
          color="green"
          :loading="isLoading"
          :disabled="isLoading"
          @click="uploadFile"
        >
          Hochladen
        </v-btn>
      </v-col>
    </v-row>
    <v-snackbar
      v-model="snackbar"
      multi-line
      location="top"
      height="96px"
      :color="snackbarStatus === 'success' ? 'green' : 'red'"
    >
      <div class="text-h6">
        {{ snackbarMsg }}
      </div>
    </v-snackbar>
  </v-container>
</template>

<script setup lang="ts">
import heic2any from "heic2any";

const state = ref("upload");
const isLoading = ref(false);
const myFile = ref();
const imageSrc = ref("");
const snackbar = ref(false);
const snackbarMsg = ref("");
const snackbarStatus = ref("success");
const btnColor = "orange";

const chooseFile = () => {
  if (document) {
    document.getElementById("fileUpload")?.click();
  }
};

const showSelectedFile = async () => {
  isLoading.value = true;
  // get image as and convert to blob
  let file = myFile.value[0];
  let blobURL = URL.createObjectURL(file);
  let blobRes = await fetch(blobURL);
  let blob = await blobRes.blob();
  // check for special ios formats and may tranform to png
  if (isHEIC(file)) {
    if (window && typeof window !== "undefined") {
      blob = (await heic2any({ blob })) as Blob;
    }
  }
  // set image view src
  const base64File = await blobToBase64(blob);
  imageSrc.value = base64File;
  state.value = "preview";
  isLoading.value = false;
};

const uploadFile = async () => {
  isLoading.value = true;
  // create file from base64String
  const base64File = imageSrc.value;
  const res: Response = await fetch(base64File);
  const blob: Blob = await res.blob();
  const file = new File([blob], "example.png", { type: "image/png" });
  // create and send http request to upload the photo
  const formData = new FormData();
  formData.append("file", file);
  const success = await $fetch<boolean>(`/upload`, {
    baseURL: "/api/v1",
    method: "POST",
    body: formData,
  });
  if (success) {
    // set image and go back
    snackbarStatus.value = "success";
    snackbarMsg.value = "Das Bild wurde erfolgreich hochgeladen";
    snackbar.value = true;
    state.value = "upload";
  } else {
    // show error msg
    snackbarStatus.value = "error";
    snackbarMsg.value = "Fehler beim Hochladen";
    snackbar.value = true;
  }
  isLoading.value = false;
};

const isHEIC = (file: File) => {
  // check file extension since windows returns blank mime for heic
  let x = file.type
    ? file.type.split("image/").pop()
    : file.name.split(".").pop()?.toLowerCase();
  return x == "heic" || x == "heif";
};

const blobToBase64 = async (blob: Blob) => {
  return new Promise<string>((resolve, _) => {
    const reader = new FileReader();
    reader.onloadend = () => resolve(reader.result as string);
    reader.readAsDataURL(blob);
  });
};
</script>

<style scoped>
.preview {
  height: calc(100vh - 320px);
}
</style>
