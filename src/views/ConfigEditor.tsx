import React, { ReactNode, useState } from 'react';
import { TabsBar, TabContent, Tab } from '@grafana/ui';
// import OAuth2Config from '../components/OAuth2Config';
import AccessTokenConfig from '../components/AccessTokenConfig';
import { ConfigEditorProps } from '../types';

interface TabListItem {
  label: string;
  key: string;
  component: (props: ConfigEditorProps) => ReactNode;
}

const tabs: TabListItem[] = [
  { label: 'Access Token', key: 'access_token', component: props => <AccessTokenConfig {...props} /> },
  // { label: 'OAuth', key: 'oauth', component: props => <OAuth2Config {...props} /> },
];

export default (props: ConfigEditorProps) => {
  const [activeTab, setActiveTab] = useState<TabListItem>(tabs[0]);

  return (
    <>
      <TabsBar>
        {tabs.map((tab, index) => {
          return (
            <Tab
              css=""
              key={index}
              label={tab.label}
              active={tab === activeTab}
              tabIndex={index}
              onChangeTab={() => setActiveTab(tab)}
            />
          );
        })}
      </TabsBar>
      <TabContent>{activeTab.component(props)}</TabContent>
    </>
  );
};
