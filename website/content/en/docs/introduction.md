---
title: Introduction
weight: 1
---

Manael is an open-source, high-performance HTTP proxy for on-the-fly image processing.

In modern web development, serving highly optimized images is crucial for improving page load speeds, reducing bandwidth costs, and boosting SEO scores. However, pre-generating and storing every permutation of an image (multiple sizes, WebP, AVIF) can quickly bloat your storage and complicate your asset management pipeline.

Manael solves this problem by processing images in real-time.

### Key Concepts

* **On-the-fly Optimization:** Manael dynamically converts original images (such as large JPEGs or PNGs) into next-generation formats like WebP and AVIF the moment they are requested.
* **Smart Content Negotiation:** By analyzing the `Accept` header of the incoming HTTP request, Manael automatically determines and serves the best image format supported by the user's browser, while seamlessly falling back to the original format for older clients.
* **Stateless by Design:** Manael does not store or cache any images itself. It acts purely as a processing layer. This makes it incredibly lightweight and easy to scale horizontally in containerized environments like Kubernetes or Cloud Run.

### Architecture

Manael is designed to sit transparently between your origin storage and your caching layer.

The optimal architecture is to place Manael behind a Content Delivery Network (CDN) and in front of your storage bucket (e.g., Amazon S3, Google Cloud Storage) or origin server.

`[Client] --> [CDN] --> [Manael Proxy] --> [Origin Storage]`

When a user requests an image, the CDN checks its cache. If it's a cache miss, the request falls back to Manael, which fetches the original image from the origin, processes it, and returns it. The CDN then caches this optimized response for all future requests.
