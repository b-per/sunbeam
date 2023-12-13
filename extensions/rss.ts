#!/usr/bin/env -S deno run -A

import Parser from "npm:rss-parser";
import { formatDistance } from "npm:date-fns";
import * as sunbeam from "https://deno.land/x/sunbeam/mod.ts";

const manifest = {
  title: "RSS",
  description: "Manage your RSS feeds",
  commands: [
    {
      name: "show",
      title: "Show a feed",
      mode: "filter",
      params: [
        {
          name: "url",
          title: "URL",
          type: "text",
        },
      ],
    },
  ],
} as const satisfies sunbeam.Manifest;

if (Deno.args.length == 0) {
  console.log(JSON.stringify(manifest));
  Deno.exit(0);
}

const payload: sunbeam.Payload<typeof manifest> = JSON.parse(Deno.args[0]);
if (payload.command == "show") {
  const feed = await new Parser().parseURL(payload.params.url);
  const page: sunbeam.List = {
    items: feed.items.map((item) => ({
      title: item.title || "",
      subtitle: item.categories?.join(", ") || "",
      accessories: item.isoDate
        ? [
            formatDistance(new Date(item.isoDate), new Date(), {
              addSuffix: true,
            }),
          ]
        : [],
      actions: [
        {
          title: "Open in browser",
          type: "open",
          url: item.link || "",
          exit: true,
        },
        {
          title: "Copy Link",
          type: "copy",
          key: "c",
          text: item.link || "",
          exit: true,
        },
      ],
    })),
  };

  console.log(JSON.stringify(page));
}
