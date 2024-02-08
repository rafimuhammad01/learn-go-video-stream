# Learn Go Video Stream

## Introduction
This repository is the repository that I used to learn about how video streaming is actually works. I code this with golang, and I am trying to design the code aligned to golang best practices or so-called idiomatic go.

## Purpose
So the purpose of this project is to create the working application to video stream. We will try to cover the streaming process from playing video and also upload a video.

Another purpose of this repository is to achieve the golang idiomatic design.

## Requirements
We will split the requirements into two big parts. One is to play the video and two is to upload the video. So the requirements for this project are:

- As a user, I want to play videos 
- As a user, I want to upload videos

## First Iteration
For the first iteration, I will implement the how to play videos first. I will have a video in my `/assets` folder. This folder will act like a cloud storage that later will be upgrade to the real cloud storage object.

In the real-world, video is being segmented to many pieces depending on how our server wants to segment it. The segmentation process will be handled later in the upload process.

Basically, the segmentation process is creating two type of files. The first one is the file that will hold the mapping for our video file. The second one is actually the video file that is already being segmented to many pieces.

You can see the file the video file example that I already generate on `/assets` file. To generate this file, I am using a segmenting format that called HLS. However, the process of segmenting file is actually having many formats, in this project we will need to handle it.

### What to build
So our purposes are actually really simple, we want to read the file directly from `/assets` folder. What makes it difficult, we need to have codes that will be easily to change the source of the file, since we will have a different source later (which is a real cloud storage object). 

The next challenge is to read the file in different format. Right now we will have the HLS format, for example, but it needs to be also able to read another format.

We also need to build our REST API to actually test our code.
