<template>
  <v-container class="pa-0 ma-0 fill-height" fluid>
    <v-row>
      <v-col>
        <v-img :src="image" class="view"></v-img>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup lang="ts">
let activeImage = 0;
const image = ref<string>();
let imageList: string[] = [];
const imageTime = 1000;
const refreshTime = 5000;

onMounted(async () => {
  // get all images and show the first one
  await listImages();
  await getImage(imageList[0]);
  // get list of images every few seconds -> there may are any new ones
  setInterval(() => {
    listImages();
  }, refreshTime);
});

const listImages = async () => {
  // get all images
  const newImageList = await $fetch<string[]>(`/listImages`, {
    baseURL: "/api/v1",
    method: "GET",
  }).catch((err) => {
    console.error(err);
    return null;
  });
  if (newImageList !== null && newImageList.length !== imageList.length) {
    // we got images and the length changed -> update
    imageList = newImageList;
  }
};

const getImage = async (imageName: string) => {
  // fetch image
  const blob = await $fetch<Blob>(`/getImage/${imageName}`, {
    baseURL: "/api/v1",
    method: "GET",
  }).catch((err: any) => {
    console.error(err);
    return null;
  });
  if (blob) {
    const base64File = await blobToBase64(blob);
    image.value = base64File;
  }
  // prepare next image
  activeImage++;
  if (activeImage >= imageList.length) {
    activeImage = 0;
  }
  setTimeout(() => getImage(imageList[activeImage]), imageTime);
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
.view {
  height: calc(100vh - 24px);
}
</style>
