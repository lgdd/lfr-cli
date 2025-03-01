/**
 * @module LiferayService
 * A service for making API calls to Liferay REST endpoints using Liferay.Util.fetch.
 * Source: https://gist.github.com/lgdd/091b8ea5952bfbc311c7febe1d5c371b
 */
const LiferayService = (() => {
  /**
   * @private
   * Makes an API call using Liferay.Util.fetch.
   *
   * @async
   * @function _fetch
   * @param {string} url The API endpoint URL.
   * @param {RequestInit} options The fetch options object.
   * @returns {Promise<any|null>} A promise that resolves with the JSON response or null (for 204).
   * @throws {Error} If the Liferay object is not defined, the response is not ok, or JSON parsing fails.
   */
  const _fetch = async (url, options) => {
    // Check if the Liferay object is defined
    if (typeof window.Liferay !== "undefined") {
      try {
        // Make the API call using Liferay's fetch method
        const response = await window.Liferay.Util.fetch(url, options);

        // Check if the response was successful
        if (!response.ok) {
          // Initialize a generic error message with HTTP status code
          let errorMessage = `An error occured (HTTP ${response.status})`;
          try {
            // Attempt to parse the error response as JSON
            const error = await response.json();
            // If parsing succeeds, append the error details
            errorMessage = `HTTP ${response.status} - ${error.title || "Unknown error"}`;
          } catch (jsonError) {
            // If parsing fails, do nothing and keep the original errorMessage
          }
          // Throw an error with the combined error message
          throw new Error(errorMessage);
        }
        // Check if the response status is 204 (No Content)
        if (response.status === 204) {
          // Return null as there is no content to parse
          return null;
        }
        try {
          // Attempt to parse the response as JSON
          return await response.json();
        } catch (jsonError) {
          // If parsing fails, log an error and throw a specific error to indicate the issue
          console.error("Error parsing JSON response:", jsonError);
          throw new Error("Error parsing JSON response");
        }
      } catch (error) {
        // Catch any errors during fetching
        console.error("Error fetching data:", error);
        // Rethrow the error to be handled by the caller
        throw error;
      }
    } else {
      // If Liferay object is not defined, log and throw an error
      console.error("Liferay (JS Object) doesn't exist.");
      throw new Error("Liferay (JS Object) doesn't exist.");
    }
  };

  return {
    /**
     * Makes a GET request to the specified URL.
     *
     * @async
     * @function get
     * @param {string} url The API endpoint URL.
     * @param {HeadersInit} [headers] Optional headers to include in the request.
     * @returns {Promise<any>} A promise that resolves with the JSON response.
     * @throws {Error} If the fetch fails or the response is not ok.
     */
    get: async (url, headers = {}) => {
      return _fetch(url, {
        method: "GET",
        headers,
      });
    },
    /**
     * Makes a POST request to the specified URL.
     *
     * @async
     * @function post
     * @param {string} url The API endpoint URL.
     * @param {any} body The request body data (will be serialized to JSON).
     * @param {HeadersInit} [headers] Optional headers to include in the request.
     * @returns {Promise<any>} A promise that resolves with the JSON response.
     * @throws {Error} If the fetch fails or the response is not ok.
     */
    post: async (url, body, headers = {}) => {
      return _fetch(url, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          ...headers,
        },
        body: JSON.stringify(body),
      });
    },
    /**
     * Makes a PUT request to the specified URL.
     *
     * @async
     * @function put
     * @param {string} url The API endpoint URL.
     * @param {any} body The request body data (will be serialized to JSON).
     * @param {HeadersInit} [headers] Optional headers to include in the request.
     * @returns {Promise<any>} A promise that resolves with the JSON response.
     * @throws {Error} If the fetch fails or the response is not ok.
     */
    put: async (url, body, headers = {}) => {
      return _fetch(url, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          ...headers,
        },
        body: JSON.stringify(body),
      });
    },
    /**
     * Makes a PATCH request to the specified URL.
     *
     * @async
     * @function patch
     * @param {string} url The API endpoint URL.
     * @param {any} body The request body data (will be serialized to JSON).
     * @param {HeadersInit} [headers] Optional headers to include in the request.
     * @returns {Promise<any>} A promise that resolves with the JSON response.
     * @throws {Error} If the fetch fails or the response is not ok.
     */
    patch: async (url, body, headers = {}) => {
      return _fetch(url, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
          ...headers,
        },
        body: JSON.stringify(body),
      });
    },
    /**
     * Makes a DELETE request to the specified URL.
     *
     * @async
     * @function delete
     * @param {string} url The API endpoint URL.
     * @param {HeadersInit} [headers] Optional headers to include in the request.
     * @returns {Promise<any>} A promise that resolves with the JSON response.
     * @throws {Error} If the fetch fails or the response is not ok.
     */
    delete: async (url, headers = {}) => {
      return _fetch(url, {
        method: "DELETE",
        headers,
      });
    },
  };
})();

export default LiferayService;