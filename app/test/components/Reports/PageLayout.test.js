import { describe, it, expect, vi, beforeEach } from "vitest";
import { mount, RouterLinkStub } from "@vue/test-utils";
import PageLayout from "~/components/reports/PageLayout.vue";

let wrapper = mount(PageLayout, {
  props: {
    currentTab: "report",
  },
  global: {
    stubs: {
      NuxtLink: RouterLinkStub,
    },
  },
});

describe("QuizAnalysisTabs.vue", () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it("renders the correct title", () => {
    const title = wrapper.find("h3");
    expect(title.text()).toBe("Analýza kvízu");
  });

  it("renders the navigation tabs", () => {
    const navItems = wrapper.findAll("li.nav-item");
    expect(navItems).toHaveLength(2);
    expect(navItems[0].text()).toBe("Otázky");
    expect(navItems[1].text()).toBe("Účastníci");
  });

  it("marks the correct tab as active based on `currentTab` prop", () => {
    const activeTab = wrapper.find(".nav-link.active");
    expect(activeTab.text()).toBe("Otázky");
  });

  it("emits `changeTab` with the correct value when a tab is clicked", async () => {
    const participantsTab = wrapper.findAll("li.nav-item")[1];
    await participantsTab.trigger("click");

    expect(wrapper.emitted().changeTab).toBeTruthy();
    expect(wrapper.emitted().changeTab[0]).toEqual(["participants"]);
  });
});
