/*
 * Virtualpaper is a service to manage users paper documents in virtual format.
 * Copyright (C) 2022  Tero Vierimaa
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import * as React from "react";
import {
  Datagrid,
  DateField,
  TextField,
  ReferenceManyField,
  Labeled,
  SimpleForm,
  Edit,
  BooleanField,
  NumberField,
  useEditController,
  CreateButton,
} from "react-admin";

import { MarkdownField } from "../Markdown";
import { useMediaQuery } from "@mui/material";

import MetadataValueCreateButton from "./ValueCreate";
import MetadataValueUpdateDialog from "./ValueEditDialog";
import { useState } from "react";

export const MetadataKeyEdit = () => {
  const { record } = useEditController();
  const [keyId, setKeyId] = useState(0);

  if (record && keyId == 0) {
    setKeyId(record.id);
  }

  const [showUpdateDialog, setShowUpdateDialog] = useState(false);
  const [valueToUpdate, setValueToUpdate] = useState({ id: 35 });

  // @ts-ignore
  const onClickValue = (id, resource, record) => {
    setValueToUpdate({
      // @ts-ignore
      record: record,
      key_id: keyId,
      id: record.id,
      basePath: "metadata/keys/" + keyId + "/values",
    });
    setShowUpdateDialog(true);
  };
  const isSmall = useMediaQuery((theme: any) => theme.breakpoints.down("sm"));

  return (
    <Edit>
      <SimpleForm>
        <MetadataValueUpdateDialog
          showDialog={showUpdateDialog}
          setShowDialog={setShowUpdateDialog}
          // @ts-ignore
          basePath={valueToUpdate.basePath}
          resource="metadata/values"
          {...valueToUpdate}
        />

        <Labeled label="Metadata key name">
          <TextField source="key" />
        </Labeled>
        <Labeled label="Description">
          <MarkdownField source="description" />
        </Labeled>

        <ReferenceManyField
          label="Values"
          reference={"metadata/values"}
          target={"key_id"}
          perPage={500}
        >
          <Datagrid
            // @ts-ignore
            rowClick={onClickValue}
            bulkActionButtons={false}
          >
            <TextField source="value" />
            <BooleanField label="Automatic matching" source="match_documents" />
            {!isSmall ? (
              <TextField label="Match by" source="match_type" />
            ) : null}
            {!isSmall ? (
              <TextField label="Filter" source="match_filter" />
            ) : null}
            <NumberField source="documents_count" label={"Total documents"} />
          </Datagrid>
        </ReferenceManyField>

        <MetadataValueCreateButton record={record} />
        <Labeled label="Created at">
          <DateField source="created_at" showTime={false} />
        </Labeled>
      </SimpleForm>
    </Edit>
  );
};