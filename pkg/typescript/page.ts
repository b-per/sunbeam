/* eslint-disable */
/**
 * This file was automatically generated by json-schema-to-typescript.
 * DO NOT MODIFY IT BY HAND. Instead, modify the source JSONSchema file,
 * and run json-schema-to-typescript to regenerate this file.
 */

export type Page = List | Detail | Form;
export type Command = Copy | Open | Run | Reload | Pop | Exit;
export type Field =
  | {
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The type of the input.
       */
      type: "text";
      /**
       * The placeholder of the input.
       */
      placeholder?: string;
      /**
       * The default value of the input.
       */
      default?: string;
      /**
       * Whether the input should be secure.
       */
      secure?: boolean;
    }
  | {
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The type of the input.
       */
      type: "checkbox";
      /**
       * The default value of the input.
       */
      default?: boolean;
      /**
       * The label of the input.
       */
      label?: string;
    }
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The type of the input.
       */
      type: "textarea";
      /**
       * The placeholder of the input.
       */
      placeholder?: string;
      /**
       * The default value of the input.
       */
      default?: string;
    }
  | {
      /**
       * The name of the input.
       */
      name: string;
      /**
       * The title of the input.
       */
      title: string;
      /**
       * The type of the input.
       */
      type: "select";
      /**
       * The items of the input.
       */
      items: {
        /**
         * The title of the item.
         */
        title: string;
        value: string | number;
      }[];
      /**
       * The default value of the input.
       */
      default?: string;
    };

export interface List {
  /**
   * The type of the page.
   */
  type: "list";
  /**
   * Whether the list should be reloaded when the query changes.
   */
  reload?: boolean;
  /**
   * The title of the list.
   */
  title?: string;
  /**
   * The items in the list.
   */
  items: Item[];
}
export interface Item {
  /**
   * The title of the item.
   */
  title: string;
  /**
   * The id of the item.
   */
  id?: string;
  /**
   * The subtitle of the item.
   */
  subtitle?: string;
  /**
   * The accessories to show on the right side of the item.
   */
  accessories?: string[];
  /**
   * The actions attached to the item.
   */
  actions?: Action[];
}
export interface Action {
  /**
   * The title of the action.
   */
  title: string;
  /**
   * The key used as a shortcut.
   */
  key?: string;
  onAction: Command;
}
export interface Copy {
  /**
   * The type of the action.
   */
  type: "copy";
  /**
   * Whether to exit the app after copying.
   */
  exit?: boolean;
  /**
   * The text to copy.
   */
  text: string;
}
export interface Open {
  /**
   * The type of the action.
   */
  type: "open";
  /**
   * Whether to exit the app after opening.
   */
  exit?: boolean;
  /**
   * The target to open.
   */
  target: string;
  app?: Application | Application[];
}
export interface Application {
  platform?: "windows" | "mac" | "linux";
  name: string;
}
export interface Run {
  /**
   * The type of the action.
   */
  type: "run";
  /**
   * The name of the command to run.
   */
  command: string;
  /**
   * The parameters to pass to the command.
   */
  params?: {
    /**
     * This interface was referenced by `undefined`'s JSON-Schema definition
     * via the `patternProperty` ".+".
     */
    [k: string]: unknown;
  };
}
export interface Reload {
  type: "reload";
  /**
   * The parameters to pass to the command.
   */
  params?: {
    /**
     * This interface was referenced by `undefined`'s JSON-Schema definition
     * via the `patternProperty` ".+".
     */
    [k: string]: unknown;
  };
}
export interface Pop {
  type: "pop";
  /**
   * Whether to reload the page after popping.
   */
  reload?: boolean;
}
export interface Exit {
  type: "exit";
}
/**
 * A detail view displaying a preview and actions.
 */
export interface Detail {
  /**
   * The type of the page.
   */
  type: "detail";
  /**
   * The title of the detail view.
   */
  title?: string;
  /**
   * The text to show in the detail view.
   */
  markdown: string;
  /**
   * The actions attached to the detail view.
   */
  actions?: Action[];
}
export interface Form {
  /**
   * The type of the page.
   */
  type: "form";
  /**
   * The title of the form.
   */
  title?: string;
  fields?: Field[];
}
